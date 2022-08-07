import Grid from "@mui/material/Grid";
import PageTitle from "../../components/parts/PageTitle";
import PageMenu from "../../components/parts/PageMenu";

const itemMenu = {
  title: "商品管理",
  icon: "coffee",
  pageInfos: [
    { title: "アイテム種別編集", url: "/admin/items/kind" },
    { title: "アイテム編集(在庫タイプ)", url: "/admin/items/stock" },
    { title: "アイテム在庫数変更", url: "/admin/items/stock/remain" },
    { title: "アイテム編集(フード)", url: "/admin/items/food" },
  ],
};

const storeMenu = {
  title: "店舗設定",
  icon: "coffee",
  pageInfos: [
    { title: "営業時間 編集", url: "/admin/store/hour" },
    { title: "特別営業時間 編集", url: "/admin/store/special_hour" },
    { title: "特別休日時間 編集", url: "/admin/store/holiday" },
  ],
};

const orderMenu = {
  title: "注文関係",
  icon: "coffee",
  pageInfos: [{ title: "予約カレンダー", url: "/admin/orders/calendar" }],
};

export default function AdminHome() {
  return (
    <>
      <PageTitle title="管理メニュー" />
      <Grid container spacing={2}>
        <Grid item xs={12} md={6}>
          <PageMenu {...itemMenu}></PageMenu>
        </Grid>
        <Grid item xs={12} md={6}>
          <PageMenu {...storeMenu}></PageMenu>
        </Grid>
        <Grid item xs={12} md={6}>
          <PageMenu {...orderMenu}></PageMenu>
        </Grid>
      </Grid>
    </>
  );
}
