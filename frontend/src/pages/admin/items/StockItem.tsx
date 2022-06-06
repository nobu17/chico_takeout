import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import StockItemTable from "./parts/StockItemTable"

export default function StockItem() {
  return (
    <>
      <AdminPage title="アイテム編集(在庫タイプ)">
        <Grid container spacing={2}>
          <Grid item xs={12}>
              <StockItemTable></StockItemTable>
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
