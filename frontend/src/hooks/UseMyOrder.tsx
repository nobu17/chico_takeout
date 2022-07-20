import { useState } from "react";
import { useAuth } from "../components/contexts/AuthContext";
import OrderApi, { UserOrderInfo } from "../libs/apis/order";

const api = new OrderApi("http://localhost:8086");

export function useMyOrder() {
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(false);
  const { state } = useAuth();
  const [activeOrder, setActiveOrder] = useState<UserOrderInfo | undefined>(
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
        setActiveOrder(res[0]);
      } else {
        setActiveOrder(undefined);
      }
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  const cancelActive = async () => {
    if (!state.isAuthorized) {
      setError(new Error("ユーザー認証がされていません。"));
      return;
    }
    if (!activeOrder || !activeOrder.id) {
      setError(new Error("有効なオーダーがありません。"));
      return;
    }
    try {
        setError(undefined);
        setLoading(true);
        await api.cancel(activeOrder.id);
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
    activeOrder,
    orderHistory,
    loadActive,
    loadHistory,
    cancelActive,
    error,
    loading,
  };
}
