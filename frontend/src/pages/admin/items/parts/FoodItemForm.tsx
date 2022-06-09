import * as React from "react";
import {
  Button,
  Container,
  MenuItem,
  Stack,
  TextField,
  FormControlLabel,
  Checkbox,
} from "@mui/material";
import { ItemKind } from "../../../../libs/ItemKind";
import { FoodItem } from "../../../../libs/FoodItem";
import { BusinessHour } from "../../../../libs/BusinessHour";
import { RhfSelects } from "../../../../components/parts/Rhf/RhfSelects";
import {
  SubmitHandler,
  useForm,
  FieldError,
  Controller,
} from "react-hook-form";

type FoodItemFormProps = {
  editItem: FoodItem;
  itemKinds: ItemKind[];
  businessHours: BusinessHour[];
  onSubmit: callbackSubmit;
  onCancel: callbackCancel;
};
interface callbackSubmit {
  (foodItem: FoodItem): void;
}
interface callbackCancel {
  (): void;
}

export type FoodItemInput = {
  id: string;
  name: string;
  description: string;
  priority: number;
  price: number;
  maxOrder: number;
  maxOrderPerDay: number;
  enabled: boolean;
  kindId: string; // for select only using kindId
  scheduleIds: string[];
};

const errorMessage = ({
  error,
  maxLength,
  min,
  max,
}: {
  name: string;
  error: FieldError | undefined;
  maxLength?: string;
  min?: string;
  max?: string;
}): string => {
  if (error?.type === "required") {
    return "入力が必要です。";
  }
  if (error?.type === "maxLength") {
    return maxLength + "文字以下にしてください。";
  }
  if (error?.type === "min") {
    return min + "以上にしてください。";
  }
  if (error?.type === "max") {
    return max + "以下にしてください。";
  }
  return "";
};

type HourSelects = {
  name: string;
  value: string;
}

const convertInput = (item: FoodItem): FoodItemInput => {
  const kindId = item.kind ? item.kind.id : "";
  return {
    id: item.id,
    name: item.name,
    description: item.description,
    priority: item.priority,
    price: item.price,
    maxOrder: item.maxOrder,
    maxOrderPerDay: item.maxOrderPerDay,
    enabled: item.enabled,
    kindId: kindId,
    scheduleIds: item.scheduleIds,
  };
};

const reverseInput = (item: FoodItemInput, kinds: ItemKind[]): FoodItem => {
  const kind = kinds.find((k) => k.id === item.kindId);
  return {
    id: item.id,
    name: item.name,
    description: item.description,
    priority: item.priority,
    price: item.price,
    maxOrder: item.maxOrder,
    maxOrderPerDay: item.maxOrderPerDay,
    enabled: item.enabled,
    kind: kind!,
    scheduleIds: item.scheduleIds,
  };
};

const convertHoursList = (hours: BusinessHour[]): HourSelects[] => {
  return hours.map(val => ({
    name: val.name,
    value: val.id
  }));
}

export default function FoodItemForm(props: FoodItemFormProps) {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<FoodItemInput>({ defaultValues: convertInput(props.editItem) });

  const onSubmit: SubmitHandler<FoodItemInput> = (data) => {
    const converted = reverseInput(data, props.itemKinds);
    props.onSubmit(converted);
  };

  const onCancel = () => {
    props.onCancel();
  };

  const hoursList = convertHoursList(props.businessHours);

  return (
    <>
      <Container maxWidth="sm" sx={{ pt: 2 }}>
        <Stack spacing={3}>
          <TextField
            label="名称"
            {...register("name", { required: true, maxLength: 15 })}
            error={Boolean(errors.name)}
            helperText={errorMessage({
              name: "名称",
              error: errors.name,
              maxLength: "15",
            })}
          />
          <TextField
            label="説明"
            multiline
            rows={5}
            {...register("description", { required: true, maxLength: 500 })}
            error={Boolean(errors.description)}
            helperText={errorMessage({
              name: "説明",
              error: errors.description,
              maxLength: "500",
            })}
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
                helperText={errorMessage({
                  name: "アイテム種別",
                  error: errors.kindId,
                })}
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
              required: true,
              min: 1,
              max: 10,
            })}
            error={Boolean(errors.priority)}
            helperText={errorMessage({
              name: "表示順序",
              error: errors.priority,
              min: "1",
              max: "15",
            })}
          />
          <TextField
            type="number"
            label="価格"
            {...register("price", {
              valueAsNumber: true,
              required: true,
              min: 1,
              max: 20000,
            })}
            error={Boolean(errors.price)}
            helperText={errorMessage({
              name: "価格",
              error: errors.price,
              min: "1",
              max: "20000",
            })}
          />
          <TextField
            type="number"
            label="最大注文可能数"
            {...register("maxOrder", {
              valueAsNumber: true,
              required: true,
              min: 1,
              max: 30,
            })}
            error={Boolean(errors.maxOrder)}
            helperText={errorMessage({
              name: "最大注文可能数",
              error: errors.maxOrder,
              min: "1",
              max: "30",
            })}
          />
          <TextField
            type="number"
            label="在庫数(日別)"
            {...register("maxOrderPerDay", {
              valueAsNumber: true,
              required: true,
              min: 1,
              max: 100,
            })}
            error={Boolean(errors.maxOrderPerDay)}
            helperText={errorMessage({
              name: "在庫数(日別)",
              error: errors.maxOrderPerDay,
              min: "1",
              max: "100",
            })}
          />
          <RhfSelects
            label="販売時間"
            name="scheduleIds"
            multiple={true}
            itemList={hoursList}
            control={control}
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
          <Button
            color="primary"
            variant="contained"
            size="large"
            onClick={handleSubmit(onSubmit)}
          >
            確定
          </Button>
          <Button
            color="error"
            variant="contained"
            size="large"
            onClick={onCancel}
          >
            キャンセル
          </Button>
        </Stack>
      </Container>
    </>
  );
}
