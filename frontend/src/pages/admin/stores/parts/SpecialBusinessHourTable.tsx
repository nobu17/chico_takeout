import * as React from "react";
import { DataGrid, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import { Button, Alert, Snackbar } from "@mui/material";

import SpecialBusinessHourFormDialog from "./SpecialBusinessHourFormDialog";
import { SpecialBusinessHour } from "../../../../libs/SpecialBusinessHour";
import useSpecialBusinessHour from "../../../../hooks/UseSpecialBusinessHour";
import useBusinessHour from "../../../../hooks/UseBusinessHour";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";
import { ApiError } from "../../../../libs/apis/apibase";

export default function SpecialBusinessHourTable() {
  const {
    specialBusinessHours,
    defaultSpecialBusinessHour,
    addSpecialBusinessHour,
    updateSpecialBusinessHour,
    deleteSpecialBusinessHour,
    loading: spHoursLoading,
    error: spHoursError,
  } = useSpecialBusinessHour();

  const {
    businessHours,
    loading: hoursLoading,
    error: hoursError,
  } = useBusinessHour();

  const [open, setOpen] = React.useState(false);
  const [item, setItem] = React.useState(defaultSpecialBusinessHour);
  const [openSnack, setOpenSnack] = React.useState(false);
  const [snackMessage, setSnackMessage] = React.useState("");
  const [uiErroMessager, setUIErrorMessage] = React.useState("");

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
    { field: "date", headerName: "日付", width: 120 },
    { field: "start", headerName: "開始時刻", width: 120 },
    { field: "end", headerName: "終了時刻", width: 120 },
    {
      field: "schedules",
      headerName: "販売時間",
      width: 190,
      valueGetter: (params) => {
        if (params.row.businessHourId) {
          return getScheduleName(params.row.businessHourId);
        }
        return "";
      },
    },
  ];

  const getScheduleName = (id: string): string => {
    const hour = businessHours.find((f) => f.id === id);
    return hour ? hour.name : "";
  };

  const handleNew = () => {
    // check max records
    setUIErrorMessage("");
    if (specialBusinessHours.length >= 10) {
      setUIErrorMessage(
        "最大10件まで作成可能です。不要データを削除して下さい。"
      );
      return;
    }
    const editItem = JSON.parse(JSON.stringify(defaultSpecialBusinessHour));
    setItem(editItem);
    setOpen(true);
  };

  const handleEdit = (item: SpecialBusinessHour) => {
    setUIErrorMessage("");
    const editItem = JSON.parse(JSON.stringify(item));
    setItem(editItem);
    setOpen(true);
  };

  const handleRemove = (item: SpecialBusinessHour) => {
    setUIErrorMessage("");
    const result = window.confirm("削除してもよろしいですか？");
    if (result) {
      deleteSpecialBusinessHour(item);
    }
  };

  const onClose = async (data: SpecialBusinessHour | null) => {
    if (data) {
      let result: ApiError | null = null;
      if (data.id !== "") {
        result = await updateSpecialBusinessHour(data);
      } else {
        result = await addSpecialBusinessHour(data);
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
      {errorMessageStr(uiErroMessager)}
      {errorMessage(spHoursError)}
      {errorMessage(hoursError)}
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
          rows={specialBusinessHours}
          columns={columns}
          disableColumnFilter={true}
          disableColumnMenu={true}
          disableColumnSelector={true}
          disableDensitySelector={true}
          editMode="row"
          hideFooter
        />
        <SpecialBusinessHourFormDialog
          open={open}
          hours={businessHours}
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
          入力に問題があります。(営業時間の重複 or 既に同一日に設定済み。)
          {snackMessage}
        </Alert>
      </Snackbar>
      <LoadingSpinner
        message="Loading..."
        isLoading={spHoursLoading || hoursLoading}
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
  },
};
