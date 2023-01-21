import OrderableInfoApi, { ItemInfo } from "../apis/orderable";
import FoodItemApi from "../apis/foodItem";
import StockItemApi from "../apis/stockItem";
import ItemKindService from "./ItemKindService";
import {
  PerDayOrderableInfo,
  OrderableItemInfo,
  CategoryItemInfo,
} from "../OrderableInfo";
import { FoodItem } from "../FoodItem";
import { StockItem } from "../StockItem";
import { AggregatedItemKind } from "../ItemKind";

export default class OrderableInfoService {
  private orderableApi: OrderableInfoApi;
  private foodApi: FoodItemApi;
  private stockApi: StockItemApi;
  private itemKindService: ItemKindService;

  constructor(public baseUrl: string = "") {
    this.orderableApi = new OrderableInfoApi(baseUrl);
    this.foodApi = new FoodItemApi(baseUrl);
    this.stockApi = new StockItemApi(baseUrl);
    this.itemKindService = new ItemKindService(baseUrl);
  }

  async get(): Promise<PerDayOrderableInfo[]> {
    const orderable = await this.orderableApi.get();
    const foods = await this.foodApi.getAll();
    const stocks = await this.stockApi.getAll();
    const kinds = await this.itemKindService.getAggregatedItemKinds();

    const lists: PerDayOrderableInfo[] = [];
    for (const perDay of orderable.data.perDayInfo) {
      const perDayInfo = {
        date: perDay.date,
        hourTypeId: perDay.hourTypeId,
        startTime: perDay.startTime,
        endTime: perDay.endTime,
        categories: this.getCategoriesItem(
          perDay.items,
          foods.data,
          stocks.data,
          kinds
        ),
      };
      lists.push(perDayInfo);
    }

    return lists;
  }

  private getCategoriesItem(
    items: ItemInfo[],
    foodItems: FoodItem[],
    stockItems: StockItem[],
    kinds: AggregatedItemKind[]
  ): CategoryItemInfo[] {
    const perKindItems: { [kindId: string]: OrderableItemInfo[] } = {};
    for (const item of items) {
      if (item.itemType === "food") {
        const foodItem = foodItems.find((f) => f.id === item.id);
        if (foodItem && foodItem.kind) {
          if (!perKindItems[foodItem.kind.id])
            perKindItems[foodItem.kind.id] = [];

          perKindItems[foodItem.kind.id].push({
            id: item.id,
            name: foodItem.name,
            type: item.itemType,
            memo: foodItem.description,
            price: foodItem.price,
            max:
              item.remain > foodItem.maxOrder ? foodItem.maxOrder : item.remain,
            imageUrl: foodItem.imageUrl,
            options: [],
          });
        }
      } else if (item.itemType === "stock") {
        const stockItem = stockItems.find((f) => f.id === item.id);
        if (stockItem && stockItem.kind) {
          if (!perKindItems[stockItem.kind.id])
            perKindItems[stockItem.kind.id] = [];

          perKindItems[stockItem.kind.id].push({
            id: item.id,
            name: stockItem.name,
            type: item.itemType,
            memo: stockItem.description,
            price: stockItem.price,
            max:
              item.remain > stockItem.maxOrder
                ? stockItem.maxOrder
                : item.remain,
            imageUrl: stockItem.imageUrl,
            options: [],
          });
        }
      }
    }

    const categories: CategoryItemInfo[] = [];
    Object.keys(perKindItems).forEach((key) => {
      const kind = kinds.find((k) => k.id === key);
      if (kind) {
        const category = {
          title: kind.name,
          kindId: kind.id,
          priority: kind.priority,
          items: perKindItems[key],
        };
        // set option item from kind
        category.items.forEach(i => i.options = kind.options)
        categories.push(category);
      }
    });
    // sort by priority
    categories.sort((a, b) => {
      if (a.priority < b.priority) return -1;
      if (a.priority > b.priority) return 1;
      return 0;
    });

    return categories;
  }
}