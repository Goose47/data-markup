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

export type BatchRq = {
  name: string;
  overlaps: number;
  priority: number;
  markups: File;
  type_id: number;
};

export type BatchCardType = {
  id: number;
  name: string;
  overlaps: number;
  priority: number;
  created_at: string;
  is_active: boolean;
  type_id: number;
};

export type AssessmentNext = {
  assessment_id: number;
  markup_type: MarkupType;
};

export type AssesmentUpdateField = {
  text: string | null;
  markup_type_field_id: number;
};

export type AssessmentUpdateRq = {
  fields: AssesmentUpdateField[];
};

export type MarkupTypeFull = RequiredFields<MarkupType, "fields">;
