import { useEffect } from "react";
import { Alert, CircularProgress, Typography } from "@mui/material";
import Grid from "@mui/material/Grid";

import PageTitle from "../../components/parts/PageTitle";
import { useMyOrder } from "../../hooks/UseMyOrder";
import ReserveInfoCard from "./parts/ReserveInfoCard";

export default function MyHome() {
  const { activeOrder, loadActive, loading, error } = useMyOrder();

  useEffect(() => {
    const init = async () => {
      await loadActive();
    };
    init();
  }, []);

  return (
    <>
      <PageTitle title="MyPage" />
      <Grid container spacing={2} direction="column" alignItems="center">
        <Grid item xs={12}>
          {loading ? (
            <CircularProgress color="success" />
          ) : error ? (
            <Alert severity="error">エラーが発生しました。{error?.message}</Alert>
          ) : (
            <ReserveInfoCard order={activeOrder}></ReserveInfoCard>
          )}
        </Grid>
      </Grid>
    </>
  );
}
