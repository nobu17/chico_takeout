import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import StockItemRemainTable from "./parts/StockItemRemainTable"

export default function StockItemRemain() {
  return (
    <>
      <AdminPage title="アイテム在庫数変更">
        <Grid container spacing={2}>
          <Grid item xs={12}>
              <StockItemRemainTable></StockItemRemainTable>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
