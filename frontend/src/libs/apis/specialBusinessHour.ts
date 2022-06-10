import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";
import { SpecialBusinessHour } from "../SpecialBusinessHour";

export default class SpecialBusinessHourApi extends ApiBase {
  async getAll(): Promise<ApiResponse<SpecialBusinessHour[]>> {
    const result = await this.getAsync<SpecialBusinessHour[]>("/store/special_hour/");
    return result;
  }
  async add(item: SpecialBusinessHour): Promise<void> {
    await this.postAsync("/store/special_hour/", item);
  }
  async update(item: SpecialBusinessHour): Promise<void> {
    await this.putAsync("/store/special_hour/" + item.id, item);
  }
  async delete(item: SpecialBusinessHour): Promise<void> {
    await this.deleteAsync("/store/special_hour/" + item.id);
  }
}
