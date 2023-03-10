import { useEffect, useState } from "react";
import {
  Alert,
  Container,
  Select,
  MenuItem,
  SelectChangeEvent,
} from "@mui/material";
import { useMessages, SelectableData } from "../../../../hooks/UseMessages";
import LoadingSpinner from "../../../../components/parts/LoadingSpinner";
import MessageForm from "./MessageForm";
import { Message } from "../../../../libs/Messages";

const errorMessage = (error: Error | undefined) => {
  if (error) {
    console.log("err", error);
    return (
      <Alert variant="filled" severity="error">
        {"エラーが発生しました。" + error}
      </Alert>
    );
  }
  return <></>;
};

export default function MessageEditContainer() {
  const { selectable, current, load, update, error, loading } = useMessages();
  const [selected, setSelected] = useState<SelectableData>();

  const handleMessageChange = async (event: SelectChangeEvent) => {
    const id = event.target.value as string;
    const foundItem = selectable.find((s) => s.id === id);
    if (foundItem) {
      setSelected(foundItem);
      await load(id);
    }
  };

  useEffect(() => {
    if (selectable && selectable.length > 0) {
      const f = async () => {
        setSelected(selectable[0]);
        await load(selectable[0].id);
      };
      f();
    }
  }, []);

  const handleSubmit = async (item: Message) => {
    await update(item);
    alert("完了しました。");
  };
  return (
    <>
      {errorMessage(error)}
      <Container maxWidth="sm" sx={{ pt: 5 }}>
        <Select
          sx={{ my: 2 }}
          fullWidth
          value={selected?.id ?? ""}
          label="メッセージ"
          onChange={handleMessageChange}
        >
          {selectable.map((item, index) => (
            <MenuItem key={index} value={item.id}>
              {item.title}
            </MenuItem>
          ))}
        </Select>
      </Container>
      <MessageForm editItem={current} onSubmit={handleSubmit}></MessageForm>
      <LoadingSpinner message="Loading..." isLoading={loading} />
    </>
  );
}
