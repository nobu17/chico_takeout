import { useState } from "react";
import { useAuth } from "../components/contexts/AuthContext";
import AdminOrderApi from "../libs/apis/adminOrder";
import OrderApi, { UserOrderInfo } from "../libs/apis/order";

const adminOrderApi = new AdminOrderApi();
const orderApi = new OrderApi();

export function useAdminOrder() {
  const [error, setError] = useState<Error>();
  const [loading, setLoading] = useState(false);
  const { state } = useAuth();
  const [orderHistory, setOrderHistory] = useState<UserOrderInfo[]>([]);

  const loadHistory = async () => {
    if (!state.isAdmin) {
      setError(new Error("管理ユーザで認証されていません。"));
      return;
    }
    try {
      setError(undefined);
      setLoading(true);
      const result = await adminOrderApi.getAll();
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

  const cancelOrder = async (orderId: string) => {
    if (!state.isAdmin) {
      setError(new Error("管理ユーザで認証されていません。"));
      return;
    }
    try {
      setError(undefined);
      setLoading(true);
      await orderApi.cancel(orderId);
      await loadHistory();
    } catch (e: any) {
      setError(e);
    } finally {
      setLoading(false);
    }
  };

  return {
    orderHistory,
    loadHistory,
    cancelOrder,
    error,
    loading,
  };
}
