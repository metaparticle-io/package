pub mod docker;
use super::Runtime;

pub trait Executor {
    fn placeholder(&self);
    fn run(&self, image: String, name: String, config: Runtime);
}