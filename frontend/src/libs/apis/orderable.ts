import ApiBase from "./apibase";
import { ApiResponse } from "./apibase";

const url = "/orderable/";

export default class OrderableInfoApi extends ApiBase {
  async get(): Promise<ApiResponse<OrderableInfo>> {
    const result = await this.getAsync<OrderableInfo>(url);
    return result;
  }
}

export type OrderableInfo = {
  startDate: string;
  endDate: string;
  perDayInfo: PerDayInfo[];
};

export type PerDayInfo = {
  date: string;
  hourTypeId: string;
  startTime: string;
  endTime: string;
  items: ItemInfo[];
};

export type ItemInfo = {
  id: string;
  itemType: string;
  remain: number;
};
