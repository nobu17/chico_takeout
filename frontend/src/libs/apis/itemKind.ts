import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";
import { ItemKind } from "../ItemKind";

export default class ItemKindApi extends ApiBase {
  async getAll(): Promise<ApiResponse<ItemKind[]>> {
    const result = await this.getAsync<ItemKind[]>("/item/kind/");
    return result;
  }

  async add(item: ItemKind): Promise<void> {
    await this.postAsync("/item/kind/", item);
  }

  async update(item: ItemKind): Promise<void> {
    await this.putAsync("/item/kind/" + item.id, item);
  }

  async delete(item: ItemKind): Promise<void> {
    await this.deleteAsync("/item/kind/" + item.id);
  }
}
