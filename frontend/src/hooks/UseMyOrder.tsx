import { useState } from "react";
import { useAuth } from "../components/contexts/AuthContext";
import OrderApi, { UserOrderInfo } from "../libs/apis/order";
import { OffsetMinutesUserCanCancel } from "../libs/Constant";
import { isBeforeFromNow } from "../libs/util/DateUtil";

const api = new OrderApi("http://localhost:8086");

export function useMyOrder() {
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(false);
  const { state } = useAuth();
  const [activeOrders, setActiveOrders] = useState<UserOrderInfo[] | undefined>(
    undefined
  );
  const [orderHistory, setOrderHistory] = useState<UserOrderInfo[]>([]);

  const loadActive = async () => {
    if (!state.isAuthorized) {
      setError(new Error("ユーザー認証がされていません。"));
      return;
    }
    try {
      setError(undefined);
      setLoading(true);
      const result = await api.getActiveByUser(state.uid);
      const res = result.data;
      if (res && res.length > 0) {
        setActiveOrders(res);
      } else {
        setActiveOrders(undefined);
      }
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const cancelActive = async (id: string) => {
    if (!state.isAuthorized) {
      setError(new Error("ユーザー認証がされていません。"));
      return;
    }
    const target = activeOrders?.find((x) => x.id === id);
    if (!target || !target.id) {
      setError(new Error("有効なオーダーがありません。"));
      return;
    }
    // not admin can not cancel before 3 hours
    if (
      !state.isAdmin &&
      !isBeforeFromNow(target.pickupDateTime, OffsetMinutesUserCanCancel)
    ) {
      alert(
        "直前の予約はできません。お手数ですが店舗にご連絡お願いいたします。"
      );
      return;
    }
    try {
      setError(undefined);
      setLoading(true);
      await api.cancel(target.id);
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
      await loadActive();
    }
  };

  const loadHistory = async () => {
    if (!state.isAuthorized) {
      setError(new Error("ユーザー認証がされていません。"));
      return;
    }
    try {
      setError(undefined);
      setLoading(true);
      const result = await api.getHistoryByUser(state.uid);
      const res = result.data;
      if (res && res.length > 0) {
        setOrderHistory(res);
      } else {
        setOrderHistory([]);
      }
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    activeOrders,
    orderHistory,
    loadActive,
    loadHistory,
    cancelActive,
    error,
    loading,
  };
}
