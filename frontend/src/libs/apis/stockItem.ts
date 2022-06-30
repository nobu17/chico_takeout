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
  
  async updateRemain(request: StockItemRemainUpdateRequest): Promise<void> {
    await this.putAsync(`/item/stock/${request.id}/remain`, request);
  }

  async delete(item: StockItem): Promise<void> {
    await this.deleteAsync("/item/stock/" + item.id);
  }
}

export type StockItemRemainUpdateRequest = {
  id: string;
  remain: number;
};

type StockItemUpdateRequest = {
  id: string;
  name: string;
  description: string;
  priority: number;
  price: number;
  maxOrder: number;
  enabled: boolean;
  imageUrl: string;
  kindId: string;
};

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
    imageUrl: item.imageUrl,
    kindId: kindId,
  };
};
