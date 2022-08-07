import ApiBase, { ApiResponse } from "./apibase";
import { UserOrderInfo } from "./order";

const getUrl = "/order/admin_all/";

export default class AdminOrderApi extends ApiBase {
  async getAll(): Promise<ApiResponse<UserOrderInfo[]>> {
    return await this.getAsync<UserOrderInfo[]>(getUrl);
  }
}
