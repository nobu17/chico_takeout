import OrderCalendar from "./OrderCalendar";
import OrderTable from "../../../../components/parts/OrderTable";
import { UserOrderInfo } from "../../../../libs/apis/order";
import { useAdminOrder } from "../../../../hooks/UseAdminOrder";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";
import { useEffect, useState } from "react";
import { Alert } from "@mui/material";
import ReserveInfoDialog from "../../../mypage/reserve/parts/ReserveInfoDialog";

export default function CalendarTable() {
  const { orderHistory, loadHistory, cancelOrder, loading, error } =
    useAdminOrder();
  const [selectedOrders, setSelectedOrders] = useState<UserOrderInfo[]>([]);
  const [open, setOpen] = useState(false);
  const [selectedItem, setSelectedItem] = useState<UserOrderInfo>();

  useEffect(() => {
    const init = async () => {
      await loadHistory();
    };
    init();
  }, []);

  const errorMessage = (error: Error | undefined) => {
    if (error) {
      console.log("err", error);
      return (
        <Alert variant="filled" severity="error">
          エラーが発生しました。
        </Alert>
      );
    }
    return <></>;
  };

  const handleCalendarSelected = (items: UserOrderInfo[]) => {
    setSelectedOrders(items);
  };

  const handleDetailSelect = (item: UserOrderInfo) => {
    setSelectedItem(item);
    setOpen(true);
  };

  const handleCancelSelected = async (item: UserOrderInfo) => {
    if (window.confirm("キャンセルしてもよろしいですか？")) {
      await cancelOrder(item.id);
    }
  };

  const onClose = () => {
    setOpen(false);
  };

  return (
    <>
      {errorMessage(error)}
      <OrderCalendar
        orders={orderHistory}
        onSelected={handleCalendarSelected}
      ></OrderCalendar>
      <div style={{ height: 600 }}>
        <OrderTable
          orders={selectedOrders}
          displays={[
            "detailButton",
            "pickupDateTime",
            "userName",
            "userEmail",
            "userTelNo",
            "total",
            "memo",
            "cancel",
            "cancelButton",
            "orderDateTime",
          ]}
          onSelected={handleDetailSelect}
          onCancelSelected={handleCancelSelected}
        ></OrderTable>
      </div>
      <ReserveInfoDialog open={open} item={selectedItem} onClose={onClose} />
      <LoadingSpinner message="Loading..." isLoading={loading} />
    </>
  );
}
