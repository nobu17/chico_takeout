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

  const checkOrderExists = async () : Promise<(Error | undefined)> => {
    if (!state.isAuthorized) {
      return new Error("ユーザー認証がされていません。");
    }
    if (state.isAdmin) {
      // admin can order unlimited
      return undefined;
    }
    try {
      setLoading(true);
      const result = await api.getActiveByUser(state.uid);
      const res = result.data;
      if (res && res.length > 0) {
        // active order exists
        return new Error("同時に予約可能な数は1つになります。再度予約したい場合は先にキャンセルをお願いします。");
      } else {
        return undefined;
      }
    } catch (e: any) {
      return new Error("エラーが発生しました。お手数ですが直接お問い合わせいただくか、再度時間をおいてご利用ください。");
    } finally {
      setLoading(false);
    }
  };

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
        userName: userInfo.name,
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
    checkOrderExists,
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
