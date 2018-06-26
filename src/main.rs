extern crate taskhero;

use taskhero::tasks::Task;

fn main() {
    println!("task: {}", Task::new("some task".to_string()));
}
