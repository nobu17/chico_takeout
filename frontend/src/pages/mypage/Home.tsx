import Grid from "@mui/material/Grid";
import PageTitle from "../../components/parts/PageTitle";

export default function MyHome() {
  return (
    <>
      <PageTitle title="MyPage" />
      <Grid container spacing={2}>
        <Grid item xs={12} md={6}></Grid>
      </Grid>
    </>
  );
}
