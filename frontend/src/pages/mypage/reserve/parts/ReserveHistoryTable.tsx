import { useEffect, useState } from "react";
import { DataGrid, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import { Button, Alert } from "@mui/material";

import LoadingSpinner from "../../../../components/parts/LoadingSpinner";
import { useMyOrder } from "../../../../hooks/UseMyOrder";
import { UserOrderInfo } from "../../../../libs/apis/order";
import ReserveInfoDialog from "./ReserveInfoDialog";

const getTotal = (order: UserOrderInfo): number => {
  const stockTotal = order.stockItems.reduce((sum, element) => {
    return sum + element.price * element.quantity;
  }, 0);
  const foodTotal = order.foodItems.reduce((sum, element) => {
    return sum + element.price * element.quantity;
  }, 0);
  return stockTotal + foodTotal;
};

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

  const columns: GridColDef[] = [
    {
      field: "id",
      width: 130,
      headerName: "",
      sortable: false,
      renderCell: (params: GridRenderCellParams<string>) => {
        return (
          <>
            <Button
              sx={{ mr: 2 }}
              variant="contained"
              onClick={(e) => handleSelect(params.row)}
            >
              詳細確認
            </Button>
          </>
        );
      },
    },
    { field: "pickupDateTime", headerName: "受取日時", width: 180 },
    { field: "orderDateTime", headerName: "注文日時", width: 180 },
    {
      field: "total",
      headerName: "合計金額",
      width: 120,
      valueGetter: (params) => {
        if (params.row) {
          return getTotal(params.row);
        }
        return "0";
      },
    },
    { field: "memo", headerName: "メモ等", width: 120 },
    {
      field: "enabled",
      headerName: "キャンセル",
      width: 120,
      valueGetter: (params) => {
        if (params.row.canceled) {
          return "○";
        }
        return "";
      },
    },
  ];

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
        <DataGrid
          sx={styles.grid}
          rows={orderHistory}
          columns={columns}
          disableColumnFilter={true}
          disableColumnMenu={true}
          disableColumnSelector={true}
          disableDensitySelector={true}
          disableSelectionOnClick={true}
          editMode="row"
          hideFooter
          getRowClassName={(params) =>
            `table-row-enabled--${params.row.canceled}`
          }
        />
      </div>
      <ReserveInfoDialog open={open} item={item} onClose={onClose} />
      <LoadingSpinner message="Loading..." isLoading={loading} />
    </>
  );
}

const styles = {
  grid: {
    ".MuiDataGrid-toolbarContainer": {
      borderBottom: "solid 1px rgba(224, 224, 224, 1)",
    },
    ".MuiDataGrid-row .MuiDataGrid-cell:not(:last-child)": {
      borderRight: "solid 1px rgba(224, 224, 224, 1) !important",
    },
    // 列ヘッダに背景色を指定
    ".MuiDataGrid-columnHeaders": {
      backgroundColor: "#65b2c6",
      color: "#fff",
    },
    // disabled row
    "& .table-row-enabled--true": {
      backgroundColor: "#696969",
      color: "#fff",
      "&:hover": {
        backgroundColor: "#696969",
        color: "#fff",
      },
    },
  },
};
