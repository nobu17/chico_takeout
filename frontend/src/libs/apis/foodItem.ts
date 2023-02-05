import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";
import { FoodItem } from "../FoodItem";

export default class FoodItemApi extends ApiBase {
  async getAll(): Promise<ApiResponse<FoodItem[]>> {
    const result = await this.getAsync<FoodItem[]>("/item/food/");
    return result;
  }

  async add(item: FoodItem): Promise<void> {
    await this.postAsync("/item/food/", convertRequest(item));
  }

  async update(item: FoodItem): Promise<void> {
    await this.putAsync("/item/food/" + item.id, convertRequest(item));
  }

  async delete(item: FoodItem): Promise<void> {
    await this.deleteAsync("/item/food/" + item.id);
  }
}

type FoodItemUpdateRequest = {
    id: string;
    name: string;
    description: string;
    priority: number;
    price: number;
    maxOrder: number;
    maxOrderPerDay: number;
    enabled: boolean;
    imageUrl: string;
    kindId: string;
    scheduleIds: string[];
    allowDates: string[];
  };
  
  const convertRequest = (item: FoodItem): FoodItemUpdateRequest => {
    const kindId = item.kind ? item.kind.id : "";
    return {
      id: item.id,
      name: item.name,
      description: item.description,
      priority: item.priority,
      price: item.price,
      maxOrder: item.maxOrder,
      maxOrderPerDay: item.maxOrderPerDay,
      enabled: item.enabled,
      imageUrl: item.imageUrl,
      kindId: kindId,
      scheduleIds: item.scheduleIds,
      allowDates: item.allowDates,
    };
  };
  