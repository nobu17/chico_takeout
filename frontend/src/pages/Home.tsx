import { useEffect } from "react";
import { Link } from "react-router-dom";
import { Box, Button, Typography, CircularProgress } from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import Grid from "@mui/material/Grid";
import { useMessages } from "../hooks/UseMessages";
import LineBreak from "../components/parts/LineBreak";

const topImage = `${process.env.PUBLIC_URL}/images/logo.jpg`;

export default function Home() {
  const { current, load, error, loading } = useMessages();

  useEffect(() => {
    const f = async () => {
      await load("1");
    };
    f();
  }, []);

  const displayMessage = () => {
    if (loading) {
      return <CircularProgress color="secondary" />;
    }

    if (error || !current) {
      // default message
      return (
        <>
          特製のスパイスカレーなどをテイクアウト予約注文できます。
          <br />
          特製のオリジナルグッズも。
          <br />
        </>
      );
    }

    return <LineBreak msg={current.content}></LineBreak>;
  };

  return (
    <>
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Box display="flex" justifyContent="center" alignItems="center">
            <img
              src={topImage}
              alt="top logo"
              style={{
                maxWidth: "100%",
                height: "auto",
                padding: "0",
                margin: "0",
              }}
            />
          </Box>
          <Typography variant="h5" align="center" gutterBottom>
            テイクアウト&イートイン 予約サイト
          </Typography>
          <Typography variant="body1" sx={{ mx: 2 }} align="center" paragraph>
            {displayMessage()}
          </Typography>
          <Grid item xs={12}>
            <Box textAlign="center">
              <Button
                variant="contained"
                color="error"
                endIcon={<SendIcon />}
                component={Link}
                to="/reserve/"
              >
                注文する
              </Button>
            </Box>
            <Box textAlign="center">
              <Typography variant="caption" align="center" gutterBottom>
                ※ユーザー登録が必要です
              </Typography>
            </Box>
          </Grid>
        </Grid>
      </Grid>
    </>
  );
}
