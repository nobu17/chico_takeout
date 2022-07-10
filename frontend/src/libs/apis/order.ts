import ApiBase, { ApiResponse } from "./apibase";

const url = "/order/";
const userActiveUrl = "/order/user/active/";

export default class OrderApi extends ApiBase {
  async add(order: OrderInfo): Promise<void> {
    return await this.postAsync(url, order);
  }
  async getActiveByUser(userId: string): Promise<ApiResponse<UserOrderInfo[]>> {
    return await this.getAsync<UserOrderInfo[]>(userActiveUrl + userId);
  }
}

export type OrderInfo = {
  userId: string;
  userName: string;
  userEmail: string;
  userTelNo: string;
  memo: string;
  pickupDateTime: string;
  stockItems: OrderItem[];
  foodItems: OrderItem[];
};

export type OrderItem = {
  itemId: string;
  quantity: number;
};


export type UserOrderInfo = {
  userId: string;
  userName: string;
  userEmail: string;
  userTelNo: string;
  memo: string;
  pickupDateTime: string;
  orderDateTime: string;
  stockItems: UserOrderItem[];
  foodItems: UserOrderItem[];
};

export type UserOrderItem = {
  itemId: string;
  name: string;
  price: number;
  quantity: number;
};
