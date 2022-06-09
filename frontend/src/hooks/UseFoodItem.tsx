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
};

const foodApi = new FoodItemApi("http://localhost:8086");
const kindApi = new ItemKindApi("http://localhost:8086");

export default function useFoodItem() {
  const [foodItems, setFoodItems] = useState<FoodItem[]>([]);
  const [itemKinds, setItemKinds] = useState<ItemKind[]>([]);
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const init = async () => {
      await getAll();
    };
    init();
  }, []);

  const addFoodItem = (item: FoodItem) => {
    add(item);
  };

  const updateFoodItem = (item: FoodItem) => {
    update(item);
  };

  const deleteFoodItem = (item: FoodItem) => {
    deletes(item);
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
    foodItems,
    itemKinds,
    defaultFoodItem,
    addFoodItem,
    updateFoodItem,
    deleteFoodItem,
    error,
    loading,
  };
}
