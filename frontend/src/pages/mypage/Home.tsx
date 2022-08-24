import { useEffect } from "react";
import { Alert, CircularProgress } from "@mui/material";
import Grid from "@mui/material/Grid";

import PageTitle from "../../components/parts/PageTitle";
import { useMyOrder } from "../../hooks/UseMyOrder";
import PageMenu from "../../components/parts/PageMenu";
import ReserveInfoCardList from "./parts/ReserveInfoCardList";

const myMenu = {
  title: "マイメニュー",
  icon: "coffee",
  pageInfos: [
    { title: "予約する", url: "/reserve/" },
    { title: "注文履歴", url: "/my_page/history/" },
  ],
};

export default function MyHome() {
  const { activeOrders, loadActive, cancelActive, loading, error } =
    useMyOrder();

  useEffect(() => {
    const init = async () => {
      await loadActive();
    };
    init();
  }, []);

  const handleCancel = async (id: string) => {
    if (window.confirm("キャンセルしてもよろしいですか？")) {
      await cancelActive(id);
    }
  };

  return (
    <>
      <PageTitle title="MyPage" />
      <Grid container spacing={2} direction="column" alignItems="center">
        <Grid item xs={12}>
          {loading ? (
            <CircularProgress color="success" />
          ) : error ? (
            <Alert severity="error">
              予約情報の取得でエラーが発生しました。{error?.message}
            </Alert>
          ) : (
            <ReserveInfoCardList
              orders={activeOrders}
              cancelRequest={handleCancel}
            ></ReserveInfoCardList>
          )}
        </Grid>
      </Grid>
      <Grid
        sx={{ my: 2 }}
        container
        spacing={2}
        direction="column"
        alignItems="center"
      >
        <Grid item xs={12} md={12}>
          <PageMenu {...myMenu}></PageMenu>
        </Grid>
      </Grid>
    </>
  );
}
