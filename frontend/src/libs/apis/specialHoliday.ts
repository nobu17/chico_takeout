import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";
import { SpecialHoliday } from "../SpecialHoliday";

const url = "/store/holiday/";

export default class SpecialHolidayApi extends ApiBase {
  async getAll(): Promise<ApiResponse<SpecialHoliday[]>> {
    const result = await this.getAsync<SpecialHoliday[]>(url);
    return result;
  }
  async add(item: SpecialHoliday): Promise<void> {
    await this.postAsync(url, item);
  }
  async update(item: SpecialHoliday): Promise<void> {
    await this.putAsync(url + item.id, item);
  }
  async delete(item: SpecialHoliday): Promise<void> {
    await this.deleteAsync(url + item.id);
  }
}
