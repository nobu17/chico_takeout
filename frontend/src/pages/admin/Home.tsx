import Grid from "@mui/material/Grid";
import PageTitle from "../../components/parts/PageTitle";
import PageMenu, { PageMenuProps } from "../../components/parts/PageMenu";

const itemMenu: PageMenuProps = {
  title: "商品管理",
  icon: "apple",
  pageInfos: [
    { title: "アイテム種別編集", url: "/admin/items/kind" },
    { title: "アイテム編集(在庫タイプ)", url: "/admin/items/stock" },
    { title: "アイテム在庫数変更", url: "/admin/items/stock/remain" },
    { title: "アイテム編集(フード)", url: "/admin/items/food" },
  ],
};

const storeMenu: PageMenuProps = {
  title: "店舗設定",
  icon: "calendar",
  pageInfos: [
    { title: "営業時間 編集", url: "/admin/store/hour" },
    { title: "特別営業時間 編集", url: "/admin/store/special_hour" },
    { title: "特別休日時間 編集", url: "/admin/store/holiday" },
  ],
};

const orderMenu: PageMenuProps = {
  title: "注文関係",
  icon: "coffee",
  pageInfos: [
    { title: "予約カレンダー", url: "/admin/orders/calendar" },
    { title: "予約一覧", url: "/admin/orders/all_orders" },
  ],
};

const statisticsMenu: PageMenuProps = {
  title: "集計",
  icon: "science",
  pageInfos: [{ title: "集計(月別)", url: "/admin/statistics/monthly" }],
};

export default function AdminHome() {
  return (
    <>
      <PageTitle title="管理メニュー" />
      <Grid container spacing={2}>
        <Grid item xs={12} md={6}>
          <PageMenu {...orderMenu}></PageMenu>
        </Grid>
        <Grid item xs={12} md={6}>
          <PageMenu {...statisticsMenu}></PageMenu>
        </Grid>
        <Grid item xs={12} md={6}>
          <PageMenu {...itemMenu}></PageMenu>
        </Grid>
        <Grid item xs={12} md={6}>
          <PageMenu {...storeMenu}></PageMenu>
        </Grid>
      </Grid>
    </>
  );
}
