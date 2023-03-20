import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";
import { BusinessHour, Schedules } from "../BusinessHour";

export default class BusinessHourApi extends ApiBase {
  async getAll(): Promise<ApiResponse<Schedules>> {
    const result = await this.getAsync<Schedules>("/store/hour/");
    return result;
  }
  async getAllEnabled(): Promise<ApiResponse<Schedules>> {
    const result = await this.getAsync<Schedules>("/store/hour/");
    result.data.schedules = result.data.schedules.filter((s) => s.enabled);
    return result;
  }
  async update(hour: BusinessHour): Promise<void> {
    await this.putAsync("/store/hour/" + hour.id, hour);
  }
  async updateEnabled(data: BusinessHourEnableUpdate): Promise<void> {
    await this.putAsync("/store/hour/" + data.id + "/enabled", data);
  }
}

type BusinessHourEnableUpdate = {
  id: string;
  enabled: boolean;
};
