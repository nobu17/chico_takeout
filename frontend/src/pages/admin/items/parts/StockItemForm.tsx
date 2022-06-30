import * as React from "react";
import {
  Container,
  MenuItem,
  Stack,
  TextField,
  FormControlLabel,
  Checkbox,
} from "@mui/material";
import { ItemKind } from "../../../../libs/ItemKind";
import { StockItem } from "../../../../libs/StockItem";
import { SubmitHandler, useForm, Controller } from "react-hook-form";
import SubmitButtons from "../../../../components/parts/SubmitButtons";
import { FirebaseImageUpload } from "../../../../components/parts/FirebaseImageUpload";

import {
  RequiredErrorMessage,
  MaxLengthErrorMessage,
  MaxErrorMessage,
  MinErrorMessage,
} from "../../../../libs/ErrorMessages";

type StockItemFormProps = {
  editItem: StockItem;
  itemKinds: ItemKind[];
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (stockItem: StockItem): void;
}
interface callbackCancel {
  (): void;
}

export type StockItemInput = {
  id: string;
  name: string;
  description: string;
  priority: number;
  price: number;
  maxOrder: number;
  enabled: boolean;
  imageUrl: string;
  remain: number;
  kindId: string; // for select only using kindId
};

const convertInput = (item: StockItem): StockItemInput => {
  const kindId = item.kind ? item.kind.id : "";
  return {
    id: item.id,
    name: item.name,
    description: item.description,
    priority: item.priority,
    price: item.price,
    maxOrder: item.maxOrder,
    enabled: item.enabled,
    imageUrl: item.imageUrl,
    remain: item.remain,
    kindId: kindId,
  };
};

const reverseInput = (item: StockItemInput, kinds: ItemKind[]): StockItem => {
  const kind = kinds.find((k) => k.id === item.kindId);
  console.log("find", kind, item, kinds);
  return {
    id: item.id,
    name: item.name,
    description: item.description,
    priority: item.priority,
    price: item.price,
    maxOrder: item.maxOrder,
    enabled: item.enabled,
    imageUrl: item.imageUrl,
    remain: item.remain,
    kind: kind!,
  };
};

const baseUrl = "/stocks";

export default function StockItemForm(props: StockItemFormProps) {
  const {
    register,
    setValue,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<StockItemInput>({ defaultValues: convertInput(props.editItem) });

  const onSubmit: SubmitHandler<StockItemInput> = (data) => {
    const converted = reverseInput(data, props.itemKinds);
    props.onSubmit(converted);
  };

  const onCancel = () => {
    props.onCancel();
  };

  const onImageChanged = (fileUrl: string) => {
    setValue("imageUrl", fileUrl);
  };

  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 2 }}>
        <Stack spacing={3}>
          <TextField
            label="名称"
            {...register("name", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 15, message: MaxLengthErrorMessage(15) },
            })}
            error={Boolean(errors.name)}
            helperText={errors.name && errors.name.message}
          />
          <TextField
            label="説明"
            multiline
            rows={5}
            {...register("description", {
              required: { value: true, message: RequiredErrorMessage },
              maxLength: { value: 500, message: MaxLengthErrorMessage(500) },
            })}
            error={Boolean(errors.description)}
            helperText={errors.description && errors.description.message}
          />
          <Controller
            name="kindId"
            control={control}
            defaultValue={""}
            rules={{ required: true }}
            render={({ field: { onChange, value, ref } }) => (
              <TextField
                inputRef={ref}
                onChange={onChange}
                value={value}
                select
                sx={{ mt: 2 }}
                fullWidth
                label="アイテム種別"
                error={Boolean(errors.kindId)}
                helperText={errors.kindId && errors.kindId.message}
              >
                {props.itemKinds.map((kind, index) => (
                  <MenuItem key={index} value={kind.id}>
                    {kind.name}
                  </MenuItem>
                ))}
              </TextField>
            )}
          />
          <TextField
            type="number"
            label="表示順序"
            {...register("priority", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 10, message: MaxErrorMessage(10) },
            })}
            error={Boolean(errors.priority)}
            helperText={errors.priority && errors.priority.message}
          />
          <TextField
            type="number"
            label="価格"
            {...register("price", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 20000, message: MaxErrorMessage(20000) },
            })}
            error={Boolean(errors.price)}
            helperText={errors.price && errors.price.message}
          />
          <TextField
            type="number"
            label="最大注文可能数"
            {...register("maxOrder", {
              valueAsNumber: true,
              required: { value: true, message: RequiredErrorMessage },
              min: { value: 1, message: MinErrorMessage(1) },
              max: { value: 30, message: MaxErrorMessage(30) },
            })}
            error={Boolean(errors.maxOrder)}
            helperText={errors.maxOrder && errors.maxOrder.message}
          />
          <Controller
            name="imageUrl"
            control={control}
            render={({ field: { onChange, value, ref } }) => (
              <FirebaseImageUpload
                baseUrl={baseUrl}
                imageUrl={value}
                onImageUploaded={onImageChanged}
              ></FirebaseImageUpload>
            )}
          />
          <FormControlLabel
            control={
              <Controller
                defaultValue={true}
                name="enabled"
                control={control}
                render={({ field: { onChange, value, ref } }) => (
                  <Checkbox
                    inputRef={ref}
                    checked={value}
                    onChange={onChange}
                  />
                )}
              />
            }
            label="有効"
          />
          <SubmitButtons
            onSubmit={handleSubmit(onSubmit)}
            onCancel={onCancel}
          />
        </Stack>
      </Container>
    </>
  );
}
