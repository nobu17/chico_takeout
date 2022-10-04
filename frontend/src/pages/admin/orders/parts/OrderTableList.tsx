import OrderTable from "../../../../components/parts/OrderTable";
import { UserOrderInfo } from "../../../../libs/apis/order";
import { useAdminOrder } from "../../../../hooks/UseAdminOrder";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";
import { useEffect, useState } from "react";
import { Alert } from "@mui/material";
import OrderDetailDialog from "../../../mypage/parts/OrderDetailDialog";

export default function OrderTableList() {
  const { orderHistory, loadHistory, cancelOrder, loading, error } =
    useAdminOrder();
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
      <div style={{ height: 600 }}>
        <OrderTable
          orders={orderHistory}
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
      <OrderDetailDialog
        open={open}
        onClose={onClose}
        order={selectedItem}
      ></OrderDetailDialog>
      <LoadingSpinner message="Loading..." isLoading={loading} />
    </>
  );
}
