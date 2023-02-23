import { useState } from "react";
import { useAuth } from "../components/contexts/AuthContext";
import OrderApi, { UserOrderInfo } from "../libs/apis/order";
import { OffsetMinutesUserCanCancel } from "../libs/Constant";
import { isBeforeFromNowStr } from "../libs/util/DateUtil";
import { useMessageDialog } from "./UseMessageDialog";

const api = new OrderApi();

export function useMyOrder() {
  const { showMessageDialog, renderDialog } = useMessageDialog();
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
    // not admin can not cancel before 1 hours
    if (
      !state.isAdmin &&
      isBeforeFromNowStr(target.pickupDateTime, OffsetMinutesUserCanCancel)
    ) {
      await showMessageDialog(
        "",
        "直前のキャンセルはできません。お手数ですが電話にてご連絡をお願いいたします。"
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
    renderDialog,
    activeOrders,
    orderHistory,
    loadActive,
    loadHistory,
    cancelActive,
    error,
    loading,
  };
}
