import ApiBase, { ApiResponse } from "./apibase";

const orderUrl = "/order/";
const userActiveUrl = "/order/user/active/";
const userHistoryUrl = "/order/user/";

export default class OrderApi extends ApiBase {
  async add(order: OrderInfo): Promise<void> {
    return await this.postAsync(orderUrl, order);
  }
  async cancel(orderId: string): Promise<void> {
    return await this.putAsync(orderUrl + orderId, {});
  }
  async getActiveByUser(userId: string): Promise<ApiResponse<UserOrderInfo[]>> {
    return await this.getAsync<UserOrderInfo[]>(userActiveUrl + userId);
  }
  async getHistoryByUser(
    userId: string
  ): Promise<ApiResponse<UserOrderInfo[]>> {
    return await this.getAsync<UserOrderInfo[]>(userHistoryUrl + userId);
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
  id: string;
  userId: string;
  userName: string;
  userEmail: string;
  userTelNo: string;
  memo: string;
  pickupDateTime: string;
  orderDateTime: string;
  stockItems: UserOrderItem[];
  foodItems: UserOrderItem[];
  canceled: boolean;
};

export type UserOrderItem = {
  itemId: string;
  name: string;
  price: number;
  quantity: number;
};

export const getTotal = (order: UserOrderInfo): number => {
  return getTotalByItems(order.stockItems, order.foodItems);
};

export const getTotalByItems = (
  stockItems: UserOrderItem[],
  foodItems: UserOrderItem[]
): number => {
  const stockTotal = stockItems.reduce((sum, element) => {
    return sum + element.price * element.quantity;
  }, 0);
  const foodTotal = foodItems.reduce((sum, element) => {
    return sum + element.price * element.quantity;
  }, 0);
  return stockTotal + foodTotal;
};
