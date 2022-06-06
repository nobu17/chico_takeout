import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";
import { StockItem } from "../StockItem";

export default class StockItemApi extends ApiBase {
  async getAll(): Promise<ApiResponse<StockItem[]>> {
    const result = await this.getAsync<StockItem[]>("/item/stock/");
    return result;
  }

  async add(item: StockItem): Promise<void> {
    await this.postAsync("/item/stock/", convertRequest(item));
  }

  async update(item: StockItem): Promise<void> {
    await this.putAsync("/item/stock/" + item.id, convertRequest(item));
  }

  async delete(item: StockItem): Promise<void> {
    await this.deleteAsync("/item/stock/" + item.id);
  }
}

type StockItemUpdateRequest = {
  id: string;
  name: string;
  description: string;
  priority: number;
  price: number;
  maxOrder: number;
  enabled: boolean;
  kindId: string;
}

const convertRequest = (item: StockItem): StockItemUpdateRequest => {
  const kindId = item.kind ? item.kind.id : "";
  return {
    id: item.id,
    name: item.name,
    description: item.description,
    priority: item.priority,
    price: item.price,
    maxOrder: item.maxOrder,
    enabled: item.enabled,
    kindId: kindId,
  };
};