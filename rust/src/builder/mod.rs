pub mod docker;

pub trait Builder {
    fn placeholder(&self);
    fn build(&self, dir: String, image: String);
    fn push(&self, image: String);
    fn logs(&self, name: String);
    fn cancel(&self, name: String);
}