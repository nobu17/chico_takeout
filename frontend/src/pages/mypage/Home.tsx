import { useEffect } from "react";
import { Alert, CircularProgress, Grid } from "@mui/material";
import PageTitle from "../../components/parts/PageTitle";
import { useMyOrder } from "../../hooks/UseMyOrder";
import PageMenu, { PageMenuProps } from "../../components/parts/PageMenu";
import ReserveInfoCardList from "./parts/ReserveInfoCardList";
import { useConfirmDialog } from "../../hooks/UseConfirmDialog";
import { useReloadTimer } from "../../hooks/UseTimer";
import { useMessages } from "../../hooks/UseMessages";
import StoreMessage from "./parts/StoreMessage";

const myMenu: PageMenuProps = {
  title: "マイメニュー",
  icon: "coffee",
  pageInfos: [
    { title: "予約する", url: "/reserve/" },
    { title: "注文履歴", url: "/my_page/history/" },
  ],
};

export default function MyHome() {
  useReloadTimer(30);
  const {
    current,
    load,
    error: messageError,
    loading: messageLoading,
  } = useMessages();
  const { showConfirmDialog, renderConfirmDialog } = useConfirmDialog();
  const {
    activeOrders,
    loadActive,
    cancelActive,
    loading,
    error,
    renderDialog,
  } = useMyOrder();

  useEffect(() => {
    const init = async () => {
      await loadActive();
      await load("2");
    };
    init();
  }, []);

  const handleCancel = async (id: string) => {
    const confirmResult = await showConfirmDialog(
      "確認",
      "キャンセルしてもよろしいですか？"
    );
    if (confirmResult) {
      await cancelActive(id);
    }
  };

  return (
    <>
      <PageTitle title="マイページ" />
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
        <Grid item xs={12}>
          {messageLoading ? (
            <CircularProgress color="success" />
          ) : messageError ? (
            <></>
          ) : (
            <StoreMessage message={current}></StoreMessage>
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
      {renderConfirmDialog()}
      {renderDialog()}
    </>
  );
}
