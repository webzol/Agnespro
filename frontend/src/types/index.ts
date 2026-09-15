// Shared frontend types (mirrors backend types in internal/types)

export interface Job {
  id: string;
  title: string;
  script: string;
  style: string;
  aspect_ratio?: string;
  episode_count?: number;
  genre_id?: string;
  genre_name?: string;
  script_model?: string;
  image_model?: string;
  video_model?: string;
  status: JobStatus;
  error?: string;
  progress: number;
  episodes: Episode[];
  characters: Character[];
  props: Prop[];
  num_episodes: number;
  num_scenes: number;
  num_characters: number;
  num_props: number;
  num_videos: number;
  created_at: string;
  updated_at: string;
}

export type JobStatus = "pending"|"planning"|"characters"|"props"|"scenes"|"video"|"done"|"failed"|"cancelled";

export interface Episode {
  id: string;
  index: number;
  title: string;
  body: string;
  summary: string;
  state: EpisodeState;
  progress: number;
  error?: string;
  characters: string[];
  props: string[];
  scenes: Scene[];
  video_url?: string;
  video_id?: string;
}

export type EpisodeState = "pending"|"running"|"characters"|"props"|"scenes"|"video"|"done"|"failed"|"skipped";

export interface Scene { id: string; index: number; heading: string; narration: string; image_prompt: string; image_url?: string; duration_s: number; }
export interface Character { id: string; job_id: string; name: string; role: string; appearance: string; image_url?: string; image_prompt?: string; }
export interface Prop { id: string; job_id: string; name: string; kind: string; description: string; image_url?: string; image_prompt?: string; }

export interface VisualStyle { id: string; name: string; category: string; desc: string; prompt_hint: string; preview: string; }
export interface Genre { id: string; name: string; desc: string; prompt_hint: string; }
export interface AspectRatio { id: string; label: string; ratio: string; image_size: string; video_size: string; hint: string; }
export interface EpisodeCountOpt { value: number; label: string; }

export interface StyleLibrary {
  visual_styles: VisualStyle[];
  genres: Genre[];
  aspect_ratios: AspectRatio[];
  episode_counts: EpisodeCountOpt[];
}

export interface ModelInfo { id: string; type: "chat"|"image"|"video"; owned_by?: string; }

export interface AdminConfig {
  api_key_set: boolean;
  api_route: string;
  base_url: string;
  script_model: string;
  image_model: string;
  video_model: string;
  concurrency: number;
  available_models: ModelInfo[];
  models_error?: string;
  fetched_at: string;
}
