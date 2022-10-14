import OrderTable from "../../../../components/parts/OrderTable";
import { UserOrderInfo } from "../../../../libs/apis/order";
import { useAdminOrder } from "../../../../hooks/UseAdminOrder";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";
import { useEffect, useState } from "react";
import { Alert } from "@mui/material";
import OrderDetailDialog from "../../../mypage/parts/OrderDetailDialog";
import OrderTableDialog from "./OrderTableDialog";

export default function OrderTableList() {
  const { orderHistory, loadHistory, cancelOrder, loading, error } =
    useAdminOrder();
  const [open, setOpen] = useState(false);
  const [openTable, setOpenTable] = useState(false);
  const [selectedItem, setSelectedItem] = useState<UserOrderInfo>();
  const [tableItems, setTableItems] = useState<UserOrderInfo[]>([]);
  const [title, setTitle] = useState("");

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

  const handleUserIdSelected = (item: UserOrderInfo) => {
    setTableItems(orderHistory.filter((o) => o.userId === item.userId));
    setTitle(`ユーザーID:${item.userId}`);
    setOpenTable(true);
  };

  const onClose = () => {
    setOpen(false);
  };

  const onTableClose = () => {
    setOpenTable(false);
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
            "userId",
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
          onUserIdSelected={handleUserIdSelected}
        ></OrderTable>
      </div>
      <OrderDetailDialog
        open={open}
        onClose={onClose}
        order={selectedItem}
      ></OrderDetailDialog>
      <OrderTableDialog
        title={title}
        open={openTable}
        onClose={onTableClose}
        orders={tableItems}
      />
      <LoadingSpinner message="Loading..." isLoading={loading} />
    </>
  );
}
