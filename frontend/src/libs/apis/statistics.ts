import ApiBase, { ApiResponse } from "./apibase";
import { MonthlyStatisticData } from "../Statistics";
import { toMonthString } from "../util/DateUtil";

const monthlyUrl = "/order/statistic/month";

export default class StatisticsApi extends ApiBase {
  async getMonthly(
    request: MonthlyRequest
  ): Promise<ApiResponse<MonthlyStatisticData>> {
    const url =
      monthlyUrl +
      `?start=${encodeURI(toMonthString(request.start))}&end=${encodeURI(
        toMonthString(request.end)
      )}`;
    return await this.getAsync<MonthlyStatisticData>(url);
  }
}

type MonthlyRequest = {
  start: Date;
  end: Date;
};
