import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";
import { BusinessHour, Schedules } from "../BusinessHour";

export default class BusinessHourApi extends ApiBase {
  async getAll(): Promise<ApiResponse<Schedules>> {
    const result = await this.getAsync<Schedules>("/store/hour/");
    return result;
  }
  async update(hour: BusinessHour) : Promise<void> {
    await this.putAsync("/store/hour/" + hour.id, hour);
  }
}
