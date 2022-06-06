import Grid from "@mui/material/Grid";
import PageTitle from "../../components/parts/PageTitle";
import PageMenu from "../../components/parts/PageMenu";

const itemMenu = {
  title: "商品管理",
  icon: "coffee",
  pageInfos: [
    { title: "アイテム種別編集", url: "/admin/items/kind" },
    { title: "link2", url: "http://aaa" },
  ],
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
          <PageMenu {...itemMenu}></PageMenu>
        </Grid>
      </Grid>
    </>
  );
}
