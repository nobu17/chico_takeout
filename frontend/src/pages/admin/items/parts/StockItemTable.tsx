import * as React from "react";
import { DataGrid, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import { Button, Alert } from "@mui/material";

import StockItemFormDialog from "./StockItemFormDialog";
import { StockItem } from "../../../../libs/StockItem";
import useStockItem from "../../../../hooks/UseStockItem";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";

export default function StockItemTable() {
  const {
    stockItems,
    defaultStockItem,
    addStockItem,
    updateStockItem,
    deleteStockItem,
    itemKinds,
    loading,
    error,
  } = useStockItem();
  const [open, setOpen] = React.useState(false);
  const [item, setItem] = React.useState(defaultStockItem);

  const columns: GridColDef[] = [
    {
      field: "id",
      width: 180,
      headerName: "",
      sortable: false,
      renderCell: (params: GridRenderCellParams<string>) => {
        return (
          <>
            <Button
              sx={{ mr: 2 }}
              variant="contained"
              onClick={(e) => handleEdit(params.row)}
            >
              編集
            </Button>
            <Button
              color="error"
              variant="contained"
              onClick={(e) => handleRemove(params.row)}
            >
              削除
            </Button>
          </>
        );
      },
    },
    { field: "priority", headerName: "表示順序", width: 100 },
    { field: "name", headerName: "アイテム名", width: 200 },
    { field: "maxOrder", headerName: "最大注文数", width: 120 },
    { field: "price", headerName: "価格(税込)", width: 120 },
    { field: "enabled", headerName: "有効", width: 120 },
    {
      field: "kind",
      headerName: "種別名",
      width: 120,
      valueGetter: (params) => {
        if (params.row.kind) {
          return params.row.kind.name;
        }
        return "";
      },
    },
  ];

  const handleRemove = (item: StockItem) => {
    const result = window.confirm("削除してもよろしいですか？");
    if (result) {
      deleteStockItem(item);
    }
  };

  const handleNew = () => {
    const editItem = JSON.parse(JSON.stringify(defaultStockItem));
    setItem(editItem);
    setOpen(true);
  };

  const handleEdit = (item: StockItem) => {
    const editItem = JSON.parse(JSON.stringify(item));
    setItem(editItem);
    setOpen(true);
  };

  const onClose = (data: StockItem | null) => {
    setOpen(false);
    if (data) {
      if (data.id === "") {
        addStockItem(data);
      } else {
        updateStockItem(data);
      }
    }
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
      <Button sx={{ my: 2 }} fullWidth variant="contained" onClick={handleNew}>
        新規作成
      </Button>
      <div style={{ height: 600 }}>
        <DataGrid
          sx={styles.grid}
          rows={stockItems}
          columns={columns}
          disableColumnFilter={true}
          disableColumnMenu={true}
          disableColumnSelector={true}
          disableDensitySelector={true}
          disableSelectionOnClick={true}
          editMode="row"
          hideFooter
          getRowClassName={(params) =>
            `table-row-enabled--${params.row.enabled}`
          }
        />
        <StockItemFormDialog
          open={open}
          editItem={item}
          itemKinds={itemKinds}
          onClose={onClose}
        />
      </div>
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
    "& .table-row-enabled--false": {
      backgroundColor: "#696969",
      color: "#fff",
      "&:hover": {
        backgroundColor: "#696969",
        color: "#fff",
      },
    },
  },
};
