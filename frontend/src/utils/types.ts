export type RequiredFields<T, K extends keyof T> = T & Required<Pick<T, K>>;

export type MarkupTypeField = {
  id?: number;
  name: string;
  label: string;
  assessment_type_id: number;
  group_id: number;
};

export type MarkupTypeRq = {
  name: string;
  fields: MarkupTypeField[];
};

export type MarkupType = {
  id: number;
  batch_id: number | null;
  name: string;
  user_id: number | null;
  fields?: MarkupTypeField[];
};

export type MarkupTypeFull = RequiredFields<MarkupType, "fields">;
