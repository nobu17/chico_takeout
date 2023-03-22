import * as React from "react";
import { DataGrid, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import { Button, Alert, Snackbar } from "@mui/material";

import {
  DayOfWeek,
  DAY_OF_WEEK,
  toShortString,
} from "../../../../libs/util/DayOfWeek";

import BusinessHourFormDialog from "./BusinessHourFormDialog";
import { BusinessHour } from "../../../../libs/BusinessHour";
import useBusinessHour from "../../../../hooks/UseBusinessHour";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";

export default function BusinessHourTable() {
  const {
    businessHours,
    defaultBusinessHour,
    updateBusinessHour,
    updateBusinessHourEnabled,
    loading,
    error,
  } = useBusinessHour();

  const [open, setOpen] = React.useState(false);
  const [item, setItem] = React.useState(defaultBusinessHour);
  const [openSnack, setOpenSnack] = React.useState(false);
  const [snackMessage, setSnackMessage] = React.useState("");

  const columns: GridColDef[] = [
    {
      field: "id",
      width: 90,
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
          </>
        );
      },
    },
    {
      field: "",
      width: 90,
      headerName: "",
      sortable: false,
      renderCell: (params: GridRenderCellParams<string>) => {
        return (
          <>
            <Button
              sx={{ mr: 2 }}
              color="error"
              variant="contained"
              onClick={(e) => handleEnableEdit(params.row, !params.row.enabled)}
            >
              {params.row.enabled ? "無効化" : "有効化"}
            </Button>
          </>
        );
      },
    },
    { field: "name", headerName: "名称", width: 120 },
    { field: "start", headerName: "開始時刻", width: 120 },
    { field: "end", headerName: "終了時刻", width: 120 },
    { field: "enabled", headerName: "有効", width: 120 },
    {
      field: "Monday",
      headerName: toShortString(DAY_OF_WEEK.Monday),
      width: 40,
      valueGetter: (params) => {
        return getDayOfWeekMark(params.row.weekdays, DAY_OF_WEEK.Monday);
      },
    },
    {
      field: "Tuesday",
      headerName: toShortString(DAY_OF_WEEK.TuesDay),
      width: 40,
      valueGetter: (params) => {
        return getDayOfWeekMark(params.row.weekdays, DAY_OF_WEEK.TuesDay);
      },
    },
    {
      field: "WednesDay",
      headerName: toShortString(DAY_OF_WEEK.WednesDay),
      width: 40,
      valueGetter: (params) => {
        return getDayOfWeekMark(params.row.weekdays, DAY_OF_WEEK.WednesDay);
      },
    },
    {
      field: "Thursday",
      headerName: toShortString(DAY_OF_WEEK.Thursday),
      width: 40,
      valueGetter: (params) => {
        return getDayOfWeekMark(params.row.weekdays, DAY_OF_WEEK.Thursday);
      },
    },
    {
      field: "Friday",
      headerName: toShortString(DAY_OF_WEEK.Friday),
      width: 40,
      valueGetter: (params) => {
        return getDayOfWeekMark(params.row.weekdays, DAY_OF_WEEK.Friday);
      },
    },
    {
      field: "Saturday",
      headerName: toShortString(DAY_OF_WEEK.Saturday),
      width: 40,
      valueGetter: (params) => {
        return getDayOfWeekMark(params.row.weekdays, DAY_OF_WEEK.Saturday);
      },
    },
    {
      field: "Sunday",
      headerName: toShortString(DAY_OF_WEEK.Sunday),
      width: 40,
      valueGetter: (params) => {
        return getDayOfWeekMark(params.row.weekdays, DAY_OF_WEEK.Sunday);
      },
    },
  ];

  const getDayOfWeekMark = (
    weekdays: number[],
    dayOfWeek: DayOfWeek
  ): string => {
    if (weekdays.includes(dayOfWeek)) {
      return "●";
    }
    return "";
  };

  const handleEdit = (item: BusinessHour) => {
    setSnackMessage("");
    setOpenSnack(false);
    const editItem = JSON.parse(JSON.stringify(item));
    setItem(editItem);
    setOpen(true);
  };

  const handleEnableEdit = async (item: BusinessHour, enabled: boolean) => {
    const msg = enabled
      ? "実行しますか？"
      : "予約等がある状態で無効化をした場合、全ての予定を無効化した場合など、想定外のエラー等が発生します。\n納得の上、自己責任で実行しますか？";
    const yesNo = window.confirm(msg);
    if (!yesNo) {
      return;
    }
    setSnackMessage("");
    setOpenSnack(false);
    const result = await updateBusinessHourEnabled({
      id: item.id,
      enabled: enabled,
    });
    if (result !== null) {
      if (result.isBadRequest()) {
        setSnackMessage(result.message);
        setOpenSnack(true);
        return;
      }
    }
  };

  const onClose = async (data: BusinessHour | null) => {
    if (data) {
      const result = await updateBusinessHour(data);
      if (result !== null) {
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

  return (
    <>
      {errorMessage(error)}
      <div style={{ height: 600 }}>
        <DataGrid
          sx={styles.grid}
          rows={businessHours}
          columns={columns}
          disableColumnFilter={true}
          disableColumnMenu={true}
          disableColumnSelector={true}
          disableDensitySelector={true}
          editMode="row"
          hideFooter
          getRowClassName={(params) =>
            `table-row-enabled--${params.row.enabled}`
          }
        />
        <BusinessHourFormDialog open={open} editItem={item} onClose={onClose} />
      </div>
      <LoadingSpinner message="Loading..." isLoading={loading} />
      <Snackbar
        open={openSnack}
        autoHideDuration={6000}
        anchorOrigin={{ vertical: "top", horizontal: "center" }}
      >
        <Alert severity="error" sx={{ width: "100%" }}>
          入力に問題があります。営業時間が重複しています。
          {snackMessage}
        </Alert>
      </Snackbar>
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
