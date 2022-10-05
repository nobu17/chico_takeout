import { useEffect, useState } from "react";
import { Alert } from "@mui/material";

import LoadingSpinner from "../../../../components/parts/LoadingSpinner";
import { useMyOrder } from "../../../../hooks/UseMyOrder";
import { UserOrderInfo } from "../../../../libs/apis/order";
import OrderTable from "../../../../components/parts/OrderTable";
import OrderDetailDialog from "../../parts/OrderDetailDialog";

export default function ReserveHistoryTable() {
  const { orderHistory, loadHistory, loading, error } = useMyOrder();
  const [open, setOpen] = useState(false);
  const [item, setItem] = useState<UserOrderInfo>();

  useEffect(() => {
    const init = async () => {
      await loadHistory();
    };
    init();
  }, []);

  const handleSelect = (item: UserOrderInfo) => {
    setItem(item);
    setOpen(true);
  };

  const onClose = () => {
    setOpen(false);
  };
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

  return (
    <>
      {errorMessage(error)}
      <div style={{ height: 600 }}>
        <OrderTable
          orders={orderHistory}
          onSelected={handleSelect}
        ></OrderTable>
      </div>
      <OrderDetailDialog
        open={open}
        onClose={onClose}
        order={item}
      ></OrderDetailDialog>
      <LoadingSpinner message="Loading..." isLoading={loading} />
    </>
  );
}
