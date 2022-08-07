import { Button } from "@mui/material";
import { DataGrid, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import { UserOrderInfo } from "../../libs/apis/order";

type OrderTableProps = {
  orders: UserOrderInfo[];
  displays?: ColumnNames[];
  onSelected: callbackSelected;
  onCancelSelected?: callbackCancelSelected;
};

interface callbackSelected {
  (item: UserOrderInfo): void;
}
interface callbackCancelSelected {
  (item: UserOrderInfo): void;
}

type ColumnNames =
  | "detailButton"
  | "pickupDateTime"
  | "orderDateTime"
  | "total"
  | "memo"
  | "cancel"
  | "cancelButton"
  | "userName"
  | "userEmail"
  | "userTelNo";

const getTotal = (order: UserOrderInfo): number => {
  const stockTotal = order.stockItems.reduce((sum, element) => {
    return sum + element.price * element.quantity;
  }, 0);
  const foodTotal = order.foodItems.reduce((sum, element) => {
    return sum + element.price * element.quantity;
  }, 0);
  return stockTotal + foodTotal;
};

export default function OrderTable(props: OrderTableProps) {
  const getColumns = (columns: ColumnNames[] | undefined): GridColDef[] => {
    const result: GridColDef[] = [];

    if (!columns) {
      columns = [
        "detailButton",
        "pickupDateTime",
        "orderDateTime",
        "total",
        "memo",
        "cancel",
      ];
    }

    for (const col of columns) {
      if (col === "detailButton") {
        result.push({
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
                  onClick={(e) => props.onSelected(params.row)}
                >
                  詳細確認
                </Button>
              </>
            );
          },
        });
        continue;
      }
      if (col === "cancelButton") {
        result.push({
          field: "cancelButton",
          width: 130,
          headerName: "",
          sortable: false,
          renderCell: (params: GridRenderCellParams<string>) => {
            if (params.row.canceled) {
              return <></>;
            }
            return (
              <>
                <Button
                  sx={{ mr: 2 }}
                  color="error"
                  variant="contained"
                  onClick={(e) => {
                    if (props.onCancelSelected) {
                      props.onCancelSelected(params.row);
                    }
                  }}
                >
                  キャンセル
                </Button>
              </>
            );
          },
        });
        continue;
      }
      if (col === "pickupDateTime") {
        result.push({
          field: "pickupDateTime",
          headerName: "受取日時",
          width: 180,
        });
        continue;
      }
      if (col === "orderDateTime") {
        result.push({
          field: "orderDateTime",
          headerName: "注文日時",
          width: 180,
        });
        continue;
      }
      if (col === "total") {
        result.push({
          field: "total",
          headerName: "合計金額",
          width: 120,
          valueGetter: (params) => {
            if (params.row) {
              return getTotal(params.row);
            }
            return "0";
          },
        });
        continue;
      }
      if (col === "memo") {
        result.push({ field: "memo", headerName: "メモ等", width: 120 });
        continue;
      }
      if (col === "cancel") {
        result.push({
          field: "enabled",
          headerName: "キャンセル",
          width: 120,
          valueGetter: (params) => {
            if (params.row.canceled) {
              return "○";
            }
            return "";
          },
        });
        continue;
      }
      if (col === "userName") {
        result.push({
          field: "userName",
          headerName: "ユーザー名",
          width: 180,
        });
        continue;
      }
      if (col === "userEmail") {
        result.push({
          field: "userEmail",
          headerName: "Email",
          width: 230,
        });
        continue;
      }
      if (col === "userTelNo") {
        result.push({
          field: "userTelNo",
          headerName: "電話番号",
          width: 180,
        });
        continue;
      }
    }

    return result;
  };
  return (
    <>
      <DataGrid
        sx={styles.grid}
        rows={props.orders}
        columns={getColumns(props.displays)}
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
