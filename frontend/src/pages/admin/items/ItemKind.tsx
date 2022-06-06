import Grid from "@mui/material/Grid";
import AdminPage from "../../../components/parts/AdminPage";
import ItemKindTable from "./parts/ItemKindTable";


export default function ItemKind() {
  return (
    <>
      <AdminPage title="アイテム種別編集">
        <Grid container spacing={2}>
          <Grid item xs={12}>
            <ItemKindTable />
          </Grid>
        </Grid>
      </AdminPage>
    </>
  );
}
