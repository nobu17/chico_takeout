import { useEffect, useState } from "react";

import { ItemKind } from "../libs/ItemKind";
import { FoodItem } from "../libs/FoodItem";
import FoodItemApi from "../libs/apis/foodItem";
import ItemKindApi from "../libs/apis/itemKind";

const defaultFoodItem: FoodItem = {
  id: "",
  name: "",
  description: "",
  priority: 1,
  price: 100,
  maxOrder: 5,
  maxOrderPerDay: 20,
  enabled: true,
  kind: null,
  scheduleIds: [],
  imageUrl: "",
  allowDates: [],
};

const allItemFilter = "全て";

const foodApi = new FoodItemApi();
const kindApi = new ItemKindApi();

export default function useFoodItem() {
  const [foodItems, setFoodItems] = useState<FoodItem[]>([]);
  const [itemKinds, setItemKinds] = useState<ItemKind[]>([]);
  const [selectedKindFilter, setSelectedKindFilter] = useState<CountDisplay>(
    CountDisplay.createAllItem()
  );
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await getAll();
    };
    init();
  }, []);

  useEffect(() => {
    updateSelectedKind();
  }, [foodItems]);

  const filteredFoods = selectedKindFilter.filterFoodItem(foodItems);
  const kindNames = kindReducer(foodItems);

  const addFoodItem = (item: FoodItem) => {
    add(item);
  };

  const updateFoodItem = (item: FoodItem) => {
    update(item);
  };

  const deleteFoodItem = (item: FoodItem) => {
    deletes(item);
  };

  const updateSelectedKind = () => {
    // if selected is not exists, reset
    if (
      !selectedKindFilter ||
      !kindNames.find((x) => selectedKindFilter?.isSameItem(x))
    ) {
      const all = kindNames.find((x) => x.isAllItem());
      if (all) {
        setSelectedKindFilter(all);
      }
    }
  };

  const updateSelectedKindFilter = (name: string) => {
    const item = kindNames.find((x) => x.isSame(name));
    if (item) {
      setSelectedKindFilter(item);
    }
  };

  const add = async (item: FoodItem) => {
    try {
      setError(undefined);
      setLoading(true);
      await foodApi.add(item);
      // reload
      await getForReload();
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const update = async (item: FoodItem) => {
    try {
      setError(undefined);
      setLoading(true);
      await foodApi.update(item);
      // reload
      await getForReload();
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const deletes = async (item: FoodItem) => {
    try {
      setError(undefined);
      setLoading(true);
      await foodApi.delete(item);
      // reload
      await getForReload();
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const getForReload = async () => {
    const result = await foodApi.getAll();
    setFoodItems(result.data);
  };

  const getAll = async () => {
    try {
      setError(undefined);
      setLoading(true);
      const result = await foodApi.getAll();
      setFoodItems(result.data);
      const kinds = await kindApi.getAll();
      setItemKinds(kinds.data);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    itemKinds,
    selectedKindFilter,
    defaultFoodItem,
    addFoodItem,
    updateFoodItem,
    deleteFoodItem,
    updateSelectedKindFilter,
    error,
    loading,
    filteredFoods,
    kindNames,
  };
}

const kindReducer = (items: FoodItem[]) => {
  const allKinds = items
    .filter((n) => n.kind != null)
    .map((n) => n.kind?.name) as string[];
  const kindCounts = allKinds.reduce<CountDisplay[]>((result, current) => {
    const element = result.find((r) => r.isSame(current));
    if (element) {
      element.increment();
    } else {
      const newItem = new CountDisplay(current, 1);
      result.push(newItem);
    }
    return result;
  }, []);

  // add all item
  kindCounts.unshift(CountDisplay.createAllItem());

  return kindCounts;
};

export class CountDisplay {
  count: number;
  name: string;
  constructor(name: string, count: number = 0) {
    this.name = name;
    this.count = count;
  }
  increment() {
    this.count++;
  }
  getName() {
    return this.name;
  }
  display() {
    // all item not display count
    if (this.isAllItem()) {
      return this.name;
    }
    return `${this.name} (${this.count})`;
  }
  isSame(name: string) {
    return this.name === name;
  }
  isSameItem(item: CountDisplay) {
    return this.name === item.name;
  }
  isAllItem() {
    return this.name === allItemFilter;
  }
  filterFoodItem(foodItems: FoodItem[]) {
    // all is not filtered
    if (this.isAllItem()) {
      return foodItems;
    }
    return foodItems.filter((x) => x.kind?.name === this.name);
  }
  static createAllItem() {
    return new CountDisplay(allItemFilter, 0);
  }
}
