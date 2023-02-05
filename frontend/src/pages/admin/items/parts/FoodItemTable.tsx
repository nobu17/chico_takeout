import * as React from "react";
import { DataGrid, GridColDef, GridRenderCellParams } from "@mui/x-data-grid";
import { Button, Alert } from "@mui/material";

import FoodItemFormDialog from "./FoodItemFormDialog";
import { FoodItem } from "../../../../libs/FoodItem";
import useFoodItem from "../../../../hooks/UseFoodItem";
import useBusinessHour from "../../../../hooks/UseBusinessHour";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";

export default function FoodItemTable() {
  const {
    foodItems,
    defaultFoodItem,
    addFoodItem,
    updateFoodItem,
    deleteFoodItem,
    itemKinds,
    loading: foodLoading,
    error: foodError,
  } = useFoodItem();

  const {
    loading: hourLoading,
    error: hourError,
    businessHours,
  } = useBusinessHour();

  const [open, setOpen] = React.useState(false);
  const [item, setItem] = React.useState(defaultFoodItem);

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
    {
      field: "kind",
      headerName: "種別",
      width: 120,
      valueGetter: (params) => {
        if (params.row.kind) {
          return params.row.kind.name;
        }
        return "";
      },
    },
    { field: "maxOrder", headerName: "最大注文数", width: 120 },
    { field: "maxOrderPerDay", headerName: "在庫数(日別)", width: 120 },
    { field: "price", headerName: "価格(税込)", width: 120 },
    { field: "enabled", headerName: "有効", width: 120 },
    {
      field: "schedules",
      headerName: "販売時間",
      width: 190,
      valueGetter: (params) => {
        if (params.row.scheduleIds) {
          return getSchedulesName(params.row.scheduleIds);
        }
        return "";
      },
    },
    {
      field: "allowDates",
      headerName: "販売期間指定",
      width: 190,
      valueGetter: (params) => {
        if (params.row.allowDates) {
          return getAllowDateNames(params.row.allowDates);
        }
        return "";
      },
    },
  ];

  const getSchedulesName = (ids: string[]): string => {
    const names: string[] = [];
    for (const id of ids) {
      for (const hour of businessHours) {
        if (hour.id === id) {
          names.push(hour.name);
        }
      }
    }
    return names.join(",");
  };

  const getAllowDateNames = (allowDates: string[]): string => {
    if (allowDates.length === 0) {
      return "指定なし";
    }
    const names: string[] = [];
    for (const allowDate of allowDates) {
      names.push(allowDate);
    }
    return names.join(",");
  };

  const handleRemove = (item: FoodItem) => {
    const result = window.confirm("削除してもよろしいですか？");
    if (result) {
      deleteFoodItem(item);
    }
  };

  const handleNew = () => {
    const editItem = JSON.parse(JSON.stringify(defaultFoodItem));
    setItem(editItem);
    setOpen(true);
  };

  const handleEdit = (item: FoodItem) => {
    const editItem = JSON.parse(JSON.stringify(item));
    setItem(editItem);
    setOpen(true);
  };

  const onClose = (data: FoodItem | null) => {
    setOpen(false);
    if (data) {
      if (data.id === "") {
        addFoodItem(data);
      } else {
        updateFoodItem(data);
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
      {errorMessage(foodError)}
      {errorMessage(hourError)}
      <Button sx={{ my: 2 }} fullWidth variant="contained" onClick={handleNew}>
        新規作成
      </Button>
      <div style={{ height: 600 }}>
        <DataGrid
          sx={styles.grid}
          rows={foodItems}
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
        <FoodItemFormDialog
          open={open}
          editItem={item}
          itemKinds={itemKinds}
          businessHours={businessHours}
          onClose={onClose}
        />
      </div>
      <LoadingSpinner
        message="Loading..."
        isLoading={foodLoading || hourLoading}
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
