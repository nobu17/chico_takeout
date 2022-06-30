import ApiBase from "./apibase";

const url = "/order/";

export default class OrderableInfoApi extends ApiBase {
  async post(order : OrderInfo): Promise<void> {
    return await this.postAsync(url, order);
  }
}

export type OrderInfo = {
    userId: string;
    memo: string;
    pickupDateTime: string;
    stockItems: OrderItem[];
    foodItems: OrderItem[];
}

export type OrderItem = {
    itemId: string;
    quantity: number;
}