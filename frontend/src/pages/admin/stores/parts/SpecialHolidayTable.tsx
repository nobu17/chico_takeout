import * as React from "react";
import { DataGrid, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import { Button, Alert, Snackbar } from "@mui/material";

import SpecialHolidayFromDialog from "./SpecialHolidayFromDialog";
import { SpecialHoliday } from "../../../../libs/SpecialHoliday";
import useSpecialHoliday from "../../../../hooks/UseSpecialHoliday";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";
import { ApiError } from "../../../../libs/apis/apibase";

export default function SpecialHolidayTable() {
  const {
    specialHolidays,
    addSpecialHoliday,
    updateSpecialHoliday,
    deleteSpecialHoliday,
    defaultSpecialHoliday,
    loading,
    error,
  } = useSpecialHoliday();

  const [open, setOpen] = React.useState(false);
  const [item, setItem] = React.useState(defaultSpecialHoliday);
  const [openSnack, setOpenSnack] = React.useState(false);
  const [snackMessage, setSnackMessage] = React.useState("");
  const [uiErrorMessage, setUiErrorMessage] = React.useState("");

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
    { field: "name", headerName: "名称", width: 120 },
    { field: "start", headerName: "開始日", width: 120 },
    { field: "end", headerName: "終了日", width: 120 },
  ];

  const handleNew = () => {
    // check max records
    setUiErrorMessage("");
    if (specialHolidays.length >= 10) {
        setUiErrorMessage(
        "最大10件まで作成可能です。不要データを削除して下さい。"
      );
      return;
    }
    const editItem = JSON.parse(JSON.stringify(defaultSpecialHoliday));
    setItem(editItem);
    setOpen(true);
  };

  const handleEdit = (item: SpecialHoliday) => {
    setUiErrorMessage("");
    const editItem = JSON.parse(JSON.stringify(item));
    setItem(editItem);
    setOpen(true);
  };

  const handleRemove = (item: SpecialHoliday) => {
    //setUIErrorMessage("");
    const result = window.confirm("削除してもよろしいですか？");
    if (result) {
      deleteSpecialHoliday(item);
    }
  };

  const onClose = async (data: SpecialHoliday | null) => {
    if (data) {
      let result: ApiError | null = null;
      if (data.id !== "") {
        result = await updateSpecialHoliday(data);
      } else {
        result = await addSpecialHoliday(data);
      }
      if (result) {
        if (result.isBadRequest()) {
          setSnackMessage(result.message);
          setOpenSnack(true);
          return;
        }
      }
    }
    setSnackMessage("");
    setOpenSnack(false);
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
  const errorMessageStr = (msg: string) => {
    if (msg) {
      return (
        <Alert variant="filled" severity="error">
          {msg}
        </Alert>
      );
    }
    return <></>;
  };

  return (
    <>
      {errorMessageStr(uiErrorMessage)}
      {errorMessage(error)}
      <div style={{ height: 600 }}>
        <Button
          sx={{ my: 2 }}
          fullWidth
          variant="contained"
          onClick={handleNew}
        >
          新規作成
        </Button>
        <DataGrid
          sx={styles.grid}
          rows={specialHolidays}
          columns={columns}
          disableColumnFilter={true}
          disableColumnMenu={true}
          disableColumnSelector={true}
          disableDensitySelector={true}
          editMode="row"
          hideFooter
        />
        <SpecialHolidayFromDialog
          open={open}
          editItem={item}
          onClose={onClose}
        />
      </div>
      <Snackbar
        open={openSnack}
        autoHideDuration={6000}
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
      >
        <Alert severity="error" sx={{ width: "100%" }}>
          入力に問題があります。(休日の重複)
          {snackMessage}
        </Alert>
      </Snackbar>
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
