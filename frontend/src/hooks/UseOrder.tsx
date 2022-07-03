import { useState } from "react";
import { useAuth } from "../components/contexts/AuthContext";
import OrderApi from "../libs/apis/order";
import { OrderInfo, OrderItem } from "../libs/apis/order";
import { Cart, ItemRequest } from "./UseItemCart";
import { PickupDate } from "./UsePickupDate";
import { UserInfo } from "./UseUserInfo";

const api = new OrderApi("http://localhost:8086");

export function useOrder() {
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(false);
  const { state } = useAuth();

  const submitOrder = async (
    pickupDate: PickupDate,
    cart: Cart,
    userInfo: UserInfo
  ): Promise<boolean> => {
    try {
      setError(undefined);
      setLoading(true);

      if (!state.isAuthorized) throw new Error("not Authorized");

      const order: OrderInfo = {
        pickupDateTime: `${pickupDate.date} ${pickupDate.time}`,
        userId: state.uid,
        userEmail: userInfo.email,
        userTelNo: userInfo.tel,
        memo: userInfo.memo,
        stockItems: getItems(cart.items, "stock"),
        foodItems: getItems(cart.items, "food"),
      };

      await api.add(order);
      return true;
    } catch (e: any) {
      setError(e);
      console.error("failed ordering", e);
      return false;
    } finally {
      setLoading(false);
    }
  };

  return {
    submitOrder,
    error,
    loading,
  };
}

const getItems = (
  items: { [id: string]: ItemRequest },
  kind: "food" | "stock"
): OrderItem[] => {
  const orders: OrderItem[] = [];
  Object.keys(items).forEach((key) => {
    if (items[key].item.type === kind) {
      orders.push({
        itemId: items[key].item.id,
        quantity: items[key].quantity,
      });
    }
  });

  return orders;
};
