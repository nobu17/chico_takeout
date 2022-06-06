import * as React from "react";
import { DataGrid, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import { Button, Alert } from "@mui/material";

import ItemKindFormDialog from "./ItemKindFormDialog";
import { ItemKind } from "../../../../libs/ItemKind";
import useItemKind from "../../../../hooks/UseItemKind";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";

export default function ItemKindTable() {
  const {
    itemKinds,
    defaultItemKind,
    addNewItemKind,
    updateItemKind,
    deleteItemKind,
    loading,
    error,
  } = useItemKind();

  const [open, setOpen] = React.useState(false);
  const [item, setItem] = React.useState(defaultItemKind);

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
    { field: "name", headerName: "アイテム名", width: 330 },
  ];

  const handleRemove = (item: ItemKind) => {
    const result = window.confirm("削除してもよろしいですか？");
    if (result) {
      deleteItemKind(item);
    }
  };

  const handleNew = () => {
    const editItem = JSON.parse(JSON.stringify(defaultItemKind));
    setItem(editItem);
    setOpen(true);
  };

  const handleEdit = (item: ItemKind) => {
    const editItem = JSON.parse(JSON.stringify(item));
    setItem(editItem);
    setOpen(true);
  };

  const onClose = (data: ItemKind | null) => {
    setOpen(false);
    if (data) {
      if (data.id === "") {
        addNewItemKind(data);
      } else {
        updateItemKind(data);
      }
    }
  };

  const errorMessage = (error: Error | undefined) => {
    if (error) {
      console.log("err", error);
      return <Alert variant="filled" severity="error">エラーが発生しました。</Alert>
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
          rows={itemKinds}
          columns={columns}
          disableColumnFilter={true}
          disableColumnMenu={true}
          disableColumnSelector={true}
          disableDensitySelector={true}
          editMode="row"
          hideFooter
        />
        <ItemKindFormDialog open={open} editItem={item} onClose={onClose} />
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
  },
};
