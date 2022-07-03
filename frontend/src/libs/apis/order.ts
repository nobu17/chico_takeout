import ApiBase from "./apibase";

const url = "/order/";

export default class OrderApi extends ApiBase {
  async add(order: OrderInfo): Promise<void> {
    return await this.postAsync(url, order);
  }
}

export type OrderInfo = {
  userId: string;
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
