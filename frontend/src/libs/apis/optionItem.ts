import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";
import { OptionItem } from "../OptionItem";

export default class OptionItemApi extends ApiBase {
  async getAll(): Promise<ApiResponse<OptionItem[]>> {
    const result = await this.getAsync<OptionItem[]>("/item/option/");
    return result;
  }

  async add(item: OptionItem): Promise<void> {
    await this.postAsync("/item/option/", item);
  }

  async update(item: OptionItem): Promise<void> {
    await this.putAsync("/item/option/" + item.id, item);
  }

  async delete(item: OptionItem): Promise<void> {
    await this.deleteAsync("/item/option/" + item.id);
  }
}
